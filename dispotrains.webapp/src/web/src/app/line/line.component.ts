import '../rxjs-operators';
import 'rxjs/add/operator/switchMap';

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params }   from '@angular/router';
import { Location }                 from '@angular/common';
import { Observable }               from 'rxjs/Observable';

import { LinesService } from '../lines.service';
import { Station, Line } from '../station';

@Component({
  selector: 'line',
  templateUrl: './line.component.html',
  styleUrls: ['./line.component.css'],
})
export class LineComponent implements OnInit {
  line: Observable<Line>;
  badStations: Observable<Station[]>;
  goodStations: Observable<Station[]>;

  constructor(
    private linesService: LinesService,
    private route: ActivatedRoute,
    private location: Location
  ) { }

  ngOnInit() {
    this.line = this.route.params
      .switchMap((params: Params) => this.linesService.getLine(params['id']));
    let stations = this.linesService.getStationsForLine(this.line);
    this.goodStations = stations.map(this.findGoodStations);
    this.badStations = stations.map(this.findBadStations);
  }

  private findGoodStations(stations: Station[]) : Station[] {
    return stations.filter((station: Station) => station.available());
  }

  private findBadStations(stations: Station[]) : Station[] {
    return stations.filter((station: Station) => !station.available());
  }
}
