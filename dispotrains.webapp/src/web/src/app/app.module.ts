import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { StationItemComponent } from './station-item/station-item.component';
import { LineListComponent } from './line-list/line-list.component';
import { LineComponent } from './line/line.component';
import { StationComponent } from './station/station.component';
import { StationStatsComponent } from './station-stats/station-stats.component';
import { ElevatorItemComponent } from './elevator-item/elevator-item.component';

@NgModule({
  declarations: [
    AppComponent,
    StationItemComponent,
    LineListComponent,
    LineComponent,
    StationComponent,
    StationStatsComponent,
    ElevatorItemComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
