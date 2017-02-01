import { Component, OnInit, OnDestroy } from '@angular/core';

import { EventService } from '../../services/event.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, OnDestroy {
    constructor(private eventService: EventService) {}

    ngOnInit() {
        this.eventService.event.subscribe((event: any) => {
            console.log(event);
        });
    }

    ngOnDestroy() {
        this.eventService.close();
    }
}
