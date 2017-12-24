import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params }   from '@angular/router';
import { Location }                 from '@angular/common';
import { Observable }               from 'rxjs/Observable';
import { combineLatest }  from 'rxjs/observable/combineLatest';

import { ElevatorStatsService } from '../elevator-stats.service';
import { StationService } from '../station.service';
import { Station, Elevator } from '../station';
import { ElevatorStats } from '../elevator-stats';

@Component({
  selector: 'app-elevator-stats',
  templateUrl: './elevator-stats.component.html',
  styleUrls: ['./elevator-stats.component.css']
})
export class ElevatorStatsComponent implements OnInit {
  stats: Observable<ElevatorStats[]>;
  station: Observable<Station>;
  elevator: Observable<Elevator>;

  constructor(
    private elevatorStatsService: ElevatorStatsService,
    private stationService: StationService,
    private route: ActivatedRoute,
    private location: Location
  ) { }

  ngOnInit() {
    this.station = this.stationService.getStation(this.route.params
      .map((params: Params) => params['id']));
    this.elevator = combineLatest(this.station, this.route.params
      .map((params: Params) => params['elevId'])).map(this._findElevator).publishLast();
    this.stats = this.elevatorStatsService.getElevatorStats(this.route.params
      .map((params: Params) => params['elevId']));
  }

  _findElevator(value: [Station, String]): Elevator {
    let station = value[0];
    let elevatorId = value[1];
    for (let elevator of station.elevators) {
      if (elevator.id === elevatorId) {
        return elevator;
      }
    }
  }
}
