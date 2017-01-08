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
  stations: Observable<Station[]>;

  constructor(
    private linesService: LinesService,
    private route: ActivatedRoute,
    private location: Location
  ) { }

  ngOnInit() {
    this.line = this.route.params
      .switchMap((params: Params) => this.linesService.getLine(params['id']));
    this.stations = this.linesService.getStationsForLine(this.line);
  }

}
