// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';
// Used by LinkButton -> Text component from grafana/ui
global.ResizeObserver = class ResizeObserver {
  //callback: ResizeObserverCallback;

  constructor(callback) {
    setTimeout(() => {
      callback(
        [
          {
            contentRect: {
              x: 1,
              y: 2,
              width: 500,
              height: 500,
              top: 100,
              bottom: 0,
              left: 100,
              right: 0,
            },
            target: {},
          },
        ],
        this
      );
    });
  }
  observe() {}
  disconnect() {}
  unobserve() {}
};

class IntersectionObserver {
  constructor(callback, options) {}
  observe() {}
  unobserve() {}
  disconnect() {}
}

Object.defineProperty(window, 'IntersectionObserver', {
  writable: true,
  configurable: true,
  value: IntersectionObserver,
});
