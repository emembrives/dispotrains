// Add the RxJS Observable operators.
import '../rxjs-operators';

import { Component, OnInit } from '@angular/core';
import { Observable }     from 'rxjs/Observable';

import { LinesService } from '../lines.service';
import { StationService } from '../station.service';
import { Line } from '../station';


@Component({
  selector: 'line-list',
  templateUrl: './line-list.component.html',
  styleUrls: ['./line-list.component.css'],
})
export class LineListComponent implements OnInit {
  lines: Promise<Line[]>;

  constructor(private linesService: LinesService) { }

  ngOnInit(): void {
    this.lines = this.linesService.getLines();
  }

  gotoDetail(line: Line): void {

  }
}
