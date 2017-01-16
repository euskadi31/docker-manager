import { Injectable } from '@angular/core';

import { Network } from '../entities/network';

import { DockerService } from './docker.service';

@Injectable()
export class NetworkService {
    constructor(private dockerService: DockerService) { }

    getNetworks(): Promise<Network[]> {
        return this.dockerService.get("/networks").then(response => response.json() as Network[]);
    }
}
