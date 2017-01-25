import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs/Rx';

import { Service } from '../entities/service';

import { DockerService } from './docker.service';

@Injectable()
export class ServiceService {
    constructor(private dockerService: DockerService) {}

    getServices(): Promise<Service[]> {
        return this.dockerService.get("/services").then(response => response.json() as Service[]);
    }
}
