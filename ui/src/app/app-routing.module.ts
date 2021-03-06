import { NgModule }              from '@angular/core';
import { RouterModule, Routes }  from '@angular/router';

import { DashboardComponent }   from './components/dashboard/dashboard.component';
import { NetworkComponent } from './components/network/network.component';
import { ServiceComponent } from './components/service/service.component';
import { NotFoundComponent } from './components/not-found/not-found.component';
import { ClusterComponent } from './components/cluster/cluster.component';
import { ImageComponent } from './components/image/image.component';
import { ServiceTaskComponent } from './components/service/service-task.component';
import { RegistryComponent } from './components/registry/registry.component';
import { ServiceLogComponent } from './components/service/service-log.component';

const appRoutes: Routes = [
    { path: 'dashboard', component: DashboardComponent },
    { path: 'network', component: NetworkComponent },
    { path: 'service', component: ServiceComponent },
    { path: 'service/:name/task', component: ServiceTaskComponent },
    { path: 'service/:name/log', component: ServiceLogComponent },
    { path: 'cluster', component: ClusterComponent },
    { path: 'image', component: ImageComponent },
    { path: 'registry', component: RegistryComponent },
    { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
    { path: '**', component: NotFoundComponent }
];

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule {}
