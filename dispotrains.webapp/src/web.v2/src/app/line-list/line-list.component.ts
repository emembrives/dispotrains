import { Component } from '@angular/core';
import { LinesService } from '../lines.service';
import { Line, NetworkStats } from '../station';

@Component({
  selector: 'app-line-list',
  templateUrl: './line-list.component.html',
  styleUrls: ['./line-list.component.css']
})
export class LineListComponent {
  lines: Promise<Line[]> | undefined;
  stats: Promise<NetworkStats> | undefined;

  constructor(private linesService: LinesService) { }

  ngOnInit() {
    this.lines = this.linesService.getLines();
    this.stats = this.linesService.getStats();
  }
}
