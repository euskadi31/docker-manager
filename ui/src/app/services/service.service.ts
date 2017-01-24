import { Injectable } from '@angular/core';

import { Service } from '../entities/service';

import { DockerService } from './docker.service';

@Injectable()
export class ServiceService {
    constructor(private dockerService: DockerService) { }

    getServices(): Promise<Service[]> {
        return this.dockerService.get("/services").then(response => response.json() as Service[]);
    }

    getLogs(name): Promise<string[]> {
        return this.dockerService.get(`/services/${name}/logs?stdout=true`).then(response => {
            return response.text().split('\n')
        });
    }

    getLogsByTimestamp(name, since): Promise<string[]> {
        return this.dockerService.get(`/services/${name}/logs?stdout=true&since=${since}`).then(response => {
            return response.text().split('\n')
        });
    }
}
