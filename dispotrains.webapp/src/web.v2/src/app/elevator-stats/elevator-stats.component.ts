import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { lastValueFrom, map } from 'rxjs';
import { ElevatorStatistics } from '../elevator-stats';
import { ElevatorStatsService } from '../elevator-stats.service';
import { Station, Elevator } from '../station';
import { StationService } from '../station.service';

@Component({
  selector: 'app-elevator-stats',
  templateUrl: './elevator-stats.component.html',
  styleUrls: ['./elevator-stats.component.css']
})
export class ElevatorStatsComponent implements OnInit {
  lineId: string | undefined
  stats: ElevatorStatistics | undefined;
  station: Station | undefined;
  elevator: Elevator | undefined;

  constructor(
    private elevatorStatsService: ElevatorStatsService,
    private stationService: StationService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.params.subscribe(async (params: Params) => {
      this.lineId = params['lineId'];
      let stationId: string = params['id'];
      this.station = await this.stationService.getStation(stationId);
      if (this.station === undefined) return;
      let elevatorId: string = params['elevId'];
      this.elevator = this._findElevator(this.station, elevatorId);
      this.stats = await this.elevatorStatsService.getElevatorStats(elevatorId);
    });
  }

  toDays(t: number): number {
    return t / (1000 * 1000 * 1000 * 3600 * 24.0);
  }

  workingRatio(es: ElevatorStatistics): number {
    return (100.0 * (es.Total - es.Broken)) / es.Total;
  }

  _findElevator(station: Station, elevatorId: String): Elevator | undefined {
    for (let elevator of station.elevators) {
      if (elevator.id === elevatorId) {
        return elevator;
      }
    }
    return undefined;
  }
}
