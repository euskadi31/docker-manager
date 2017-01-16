import { Injectable } from '@angular/core';

import { Task } from '../entities/task';

import { DockerService } from './docker.service';

@Injectable()
export class TaskService {
    constructor(private dockerService: DockerService) { }

    getTasks(service_name: string): Promise<Task[]> {
        return this.dockerService.get(`/tasks?service=${service_name}`)
            .then(response => response.json() as Task[])
            .then(tasks => {
                return tasks.map(task => {
                    task.Name = `${service_name}.${task.Slot}`;
                    task.ContainerID = `${task.Name}.${task.ID}`

                    return task;
                })
            });
    }
}
