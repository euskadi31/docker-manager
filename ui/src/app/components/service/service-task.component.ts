import { Component, OnInit } from '@angular/core';

import { Task } from '../../entities/task';

import { TaskService } from '../../services/task.service';

@Component({
    selector: 'app-service-task',
    templateUrl: './service-task.component.html',
    styleUrls: ['./service-task.component.css']
})
export class ServiceTaskComponent implements OnInit {
    tasks: Task[];

    constructor(private taskService: TaskService) {
        this.tasks = [];
    }

    ngOnInit() {
        this.taskService.getTasks('app').then(tasks => {
            console.log(tasks);
            this.tasks = tasks;
        });
    }

}
