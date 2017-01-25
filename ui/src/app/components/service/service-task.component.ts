import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, Params} from '@angular/router';

import { Task } from '../../entities/task';

import { TaskService } from '../../services/task.service';

@Component({
    selector: 'app-service-task',
    templateUrl: './service-task.component.html',
    styleUrls: ['./service-task.component.css']
})
export class ServiceTaskComponent implements OnInit {
    tasks: Task[];

    constructor(
        private taskService: TaskService,
        private activatedRoute: ActivatedRoute
    ) {
        this.tasks = [];
    }

    fetchTask(name) {
        this.taskService.getTasks(name).then(tasks => {
            console.log(tasks);
            this.tasks = tasks;
        });
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            this.fetchTask(params['name']);
        });
    }

}
