import { Injectable } from '@angular/core';

import { Swarm } from '../entities/swarm';

import { DockerService } from './docker.service';

@Injectable()
export class SwarmService {
    constructor(private dockerService: DockerService) { }

    getInfo(): Promise<Swarm> {
        return this.dockerService.get("/swarm").then(response => response.json() as Swarm);
    }
}
