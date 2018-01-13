// Add the RxJS Observable operators.
import '../rxjs-operators';

import { Component, OnInit } from '@angular/core';
import { Observable }     from 'rxjs/Observable';

import { LinesService } from '../lines.service';
import { StationService } from '../station.service';
import { Line, NetworkStats } from '../station';


@Component({
  selector: 'line-list',
  templateUrl: './line-list.component.html',
  styleUrls: ['./line-list.component.css'],
})
export class LineListComponent implements OnInit {
  lines: Promise<Line[]>;
  stats: Promise<NetworkStats>;

  constructor(private linesService: LinesService) { }

  ngOnInit(): void {
    this.lines = this.linesService.getLines();
    this.stats = this.linesService.getStats();
  }
}
