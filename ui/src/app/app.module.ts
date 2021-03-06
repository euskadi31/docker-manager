import { BrowserModule } from '@angular/platform-browser';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './components/app/app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { NetworkComponent } from './components/network/network.component';
import { ServiceComponent } from './components/service/service.component';
import { ServiceTaskComponent } from './components/service/service-task.component';
import { NotFoundComponent } from './components/not-found/not-found.component';
import { PanelComponent } from './components/panel/panel.component';
import { ClusterComponent } from './components/cluster/cluster.component';
import { ImageComponent } from './components/image/image.component';
import { RegistryComponent } from './components/registry/registry.component';
import { ServiceLogComponent } from './components/service/service-log.component';

import { NetworkService } from './services/network.service';
import { DockerService } from './services/docker.service';
import { ServiceService } from './services/service.service';
import { ImageService } from './services/image.service';
import { SwarmService } from './services/swarm.service';
import { NodeService } from './services/node.service';
import { TaskService } from './services/task.service';
import { RegistryService } from './services/registry.service';
import { WebSocketService } from './services/websocket.service';
import { EventService } from './services/event.service';

import { TruncatePipe } from './pipes/truncate.pipe';
import { SizePipe } from './pipes/size.pipe';
import { DialogComponent } from './components/dialog/dialog.component';
import { ConsoleComponent } from './components/console/console.component';



@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        NetworkComponent,
        ServiceComponent,
        NotFoundComponent,
        TruncatePipe,
        PanelComponent,
        ClusterComponent,
        ImageComponent,
        ServiceTaskComponent,
        SizePipe,
        RegistryComponent,
        DialogComponent,
        ServiceLogComponent,
        ConsoleComponent
    ],
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        AppRoutingModule
    ],
    providers: [
        NetworkService,
        DockerService,
        ServiceService,
        NodeService,
        ImageService,
        SwarmService,
        TaskService,
        RegistryService,
        WebSocketService,
        EventService
    ],
    schemas: [ NO_ERRORS_SCHEMA ],
    bootstrap: [AppComponent]
})
export class AppModule { }
