import { Component, OnInit, OnDestroy } from '@angular/core';

import { WebSocketService } from '../../services/websocket.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, OnDestroy {
    constructor(private websocketService: WebSocketService) {}

    ngOnInit() {
        const protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

        const url = `${protocol}//${location.host}/ws/events`;

        this.websocketService.connect(url);

        this.websocketService.message$.subscribe((event: MessageEvent) => {
            let item = JSON.parse(event.data);

            console.log(item);
        });
    }

    ngOnDestroy() {
        this.websocketService.close();
    }
}
