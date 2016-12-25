export class SorterUtils {
  static sorterBySelector(selector: (obj: any) => any) {
    return function(a: any, b: any) {
        if (selector(a) > selector(b)) {
          return 1;
        }

        if (selector(a) < selector(b)) {
          return -1;
        }
        return 0;
      };
  }
}
