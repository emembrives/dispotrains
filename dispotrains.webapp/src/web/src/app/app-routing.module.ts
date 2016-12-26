import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { LineListComponent } from './line-list/line-list.component';
import { LineComponent } from './line/line.component';
import { StationComponent } from './station/station.component';
import { StationStatsComponent } from './station-stats/station-stats.component';

const routes: Routes = [
  { path: '', redirectTo: '/lignes', pathMatch: 'full' },
  { path: 'lignes', component: LineListComponent },
  { path: 'ligne/:id', component: LineComponent },
  { path: 'gare/:id', component: StationComponent },
  { path: 'gare/:id/stats', component: StationStatsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  providers: []
})
export class AppRoutingModule { }
