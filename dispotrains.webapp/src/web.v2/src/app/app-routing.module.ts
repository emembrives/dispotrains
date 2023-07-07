import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AboutComponent } from './about/about.component';
import { ElevatorStatsComponent } from './elevator-stats/elevator-stats.component';
import { LineListComponent } from './line-list/line-list.component';
import { LineComponent } from './line/line.component';
import { StationComponent } from './station/station.component';

const routes: Routes = [
  { path: '', redirectTo: '/lignes', pathMatch: 'full' },
  { path: 'about', component: AboutComponent },
  { path: 'lignes', component: LineListComponent },
  { path: 'ligne/:id', component: LineComponent },
  { path: 'ligne/:lineId/:id', component: StationComponent },
  { path: 'ligne/:lineId/:id/:elevId', component: ElevatorStatsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
