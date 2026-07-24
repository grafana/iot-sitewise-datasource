// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';
// Used by LinkButton -> Text component from grafana/ui
import { ReadableStream, TransformStream, WritableStream } from 'stream/web';
import { TextDecoder, TextEncoder } from 'util';
import { MessageChannel, MessagePort } from 'worker_threads';

global.TextEncoder = TextEncoder;
global.TextDecoder = TextDecoder;
global.ReadableStream = ReadableStream;
global.TransformStream = TransformStream;
global.WritableStream = WritableStream;
global.MessageChannel = MessageChannel;
global.MessagePort = MessagePort;

const mockIntersectionObserver = jest.fn().mockImplementation((arg) => ({
  observe: jest.fn().mockImplementation((elem) => {
    arg([{ target: elem, isIntersecting: true }]);
  }),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

global.IntersectionObserver = mockIntersectionObserver;

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
