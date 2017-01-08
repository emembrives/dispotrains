import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { AppRoutingModule } from './app-routing.module';
import { MaterialModule } from '@angular/material';

import { AppComponent } from './app.component';
import { StationItemComponent } from './station-item/station-item.component';
import { LineListComponent } from './line-list/line-list.component';
import { LineComponent } from './line/line.component';
import { StationComponent } from './station/station.component';
import { StationStatsComponent } from './station-stats/station-stats.component';
import { ElevatorItemComponent } from './elevator-item/elevator-item.component';
import { TitlebarComponent } from './titlebar/titlebar.component';

import { StationService } from './station.service';
import { LinesService } from './lines.service';

@NgModule({
  declarations: [
    AppComponent,
    StationItemComponent,
    LineListComponent,
    LineComponent,
    StationComponent,
    StationStatsComponent,
    ElevatorItemComponent,
    TitlebarComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule,
    MaterialModule.forRoot()
  ],
  providers: [StationService, LinesService],
  bootstrap: [AppComponent]
})
export class AppModule { }
