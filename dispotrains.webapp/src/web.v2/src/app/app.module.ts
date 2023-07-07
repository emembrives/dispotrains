import { NgModule, isDevMode } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TitlebarComponent } from './titlebar/titlebar.component';
import { AboutComponent } from './about/about.component';
import { ElevatorStatsComponent } from './elevator-stats/elevator-stats.component';
import { LineComponent } from './line/line.component';
import { LineListComponent } from './line-list/line-list.component';
import { StationComponent } from './station/station.component';
import { StationItemComponent } from './station-item/station-item.component';
import { ElevatorStatsService } from './elevator-stats.service';
import { LinesService } from './lines.service';
import { StationService } from './station.service';
import { HttpClientModule } from '@angular/common/http';
import { ServiceWorkerModule } from '@angular/service-worker';

@NgModule({
  declarations: [
    AppComponent,
    TitlebarComponent,
    AboutComponent,
    ElevatorStatsComponent,
    LineComponent,
    LineListComponent,
    StationComponent,
    StationItemComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    ServiceWorkerModule.register('ngsw-worker.js', {
      enabled: !isDevMode(),
      // Register the ServiceWorker as soon as the application is stable
      // or after 30 seconds (whichever comes first).
      registrationStrategy: 'registerWhenStable:30000'
    }),
  ],
  providers: [StationService, LinesService, ElevatorStatsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
