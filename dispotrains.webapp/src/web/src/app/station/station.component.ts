import '../rxjs-operators';
import 'rxjs/add/operator/switchMap';

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params }   from '@angular/router';
import { Location }                 from '@angular/common';
import { Observable }               from 'rxjs/Observable';

import { StationService } from '../station.service';
import { Station } from '../station';

@Component({
  selector: 'station',
  templateUrl: './station.component.html',
  styleUrls: ['./station.component.css']
})
export class StationComponent implements OnInit {
  station: Observable<Station>;

  constructor(
    private stationService: StationService,
    private route: ActivatedRoute,
    private location: Location
  ) { }

  ngOnInit() {
    this.station = this.stationService.getStation(this.route.params
      .map((params: Params) => params['id']));
  }
}
