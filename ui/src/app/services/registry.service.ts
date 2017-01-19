import { Injectable } from '@angular/core';

import { Registry } from '../entities/registry';

import { DockerService } from './docker.service';

@Injectable()
export class RegistryService {
    constructor(private dockerService: DockerService) { }

    getRegistries(): Promise<Registry[]> {
        return this.dockerService.get("/registries").then(response => response.json() as Registry[]);
    }

    addRegistry(registry: Registry): Promise<Registry> {
        return this.dockerService.post("/registries", registry).then(response => response.json() as Registry);
    }
}
