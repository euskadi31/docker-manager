import { Injectable } from '@angular/core';

import { Node } from '../entities/node';

import { DockerService } from './docker.service';

@Injectable()
export class NodeService {
    constructor(private dockerService: DockerService) { }

    getNodes(): Promise<Node[]> {
        return this.dockerService.get("/nodes").then(response => response.json() as Node[]);
    }
}
