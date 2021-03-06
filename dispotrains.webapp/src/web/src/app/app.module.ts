import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { MatListModule, MatIconModule, MatToolbarModule, MatCardModule } from '@angular/material';
import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { StationItemComponent } from './station-item/station-item.component';
import { LineListComponent } from './line-list/line-list.component';
import { LineComponent } from './line/line.component';
import { StationComponent } from './station/station.component';
import { ElevatorStatsComponent } from './elevator-stats/elevator-stats.component';
import { ElevatorItemComponent } from './elevator-item/elevator-item.component';
import { TitlebarComponent } from './titlebar/titlebar.component';

import { ElevatorStatsService } from './elevator-stats.service';
import { StationService } from './station.service';
import { LinesService } from './lines.service';
import { AboutComponent } from './about/about.component';

@NgModule({
  declarations: [
    AboutComponent,
    AppComponent,
    StationItemComponent,
    LineListComponent,
    LineComponent,
    StationComponent,
    ElevatorStatsComponent,
    ElevatorItemComponent,
    TitlebarComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule,
    MatListModule,
    MatIconModule,
    MatToolbarModule,
    MatCardModule,
  ],
  providers: [StationService, LinesService, ElevatorStatsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
