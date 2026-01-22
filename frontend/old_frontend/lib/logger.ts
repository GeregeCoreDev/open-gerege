import * as Sentry from '@sentry/nextjs';

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

export interface LogEntry {
  level: LogLevel;
  message: string;
  context?: Record<string, unknown>;
  timestamp: Date;
}

const isDev = process.env.NODE_ENV === 'development';

export const logger = {
  debug: (message: string, context?: Record<string, unknown>) => log('debug', message, context),
  info: (message: string, context?: Record<string, unknown>) => log('info', message, context),
  warn: (message: string, context?: Record<string, unknown>) => log('warn', message, context),
  error: (message: string, context?: Record<string, unknown>) => log('error', message, context),
};

export const createLogger = (namespace: string) => ({
  debug: (message: string, context?: Record<string, unknown>) =>
    logger.debug(`[${namespace}] ${message}`, context),
  info: (message: string, context?: Record<string, unknown>) =>
    logger.info(`[${namespace}] ${message}`, context),
  warn: (message: string, context?: Record<string, unknown>) =>
    logger.warn(`[${namespace}] ${message}`, context),
  error: (message: string, context?: Record<string, unknown>) =>
    logger.error(`[${namespace}] ${message}`, context),
});

function log(level: LogLevel, message: string, context?: Record<string, unknown>) {
  const entry: LogEntry = {
    level,
    message,
    context,
    timestamp: new Date(),
  };

  if (isDev) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const args: any[] = [message];
    if (context) args.push(context);


    console[level](...args);
  } else {
    // Send to Sentry
    if (level === 'error') {
      Sentry.captureException(new Error(message), {
        extra: context,
        level: 'error',
      });

      console.error(JSON.stringify(entry));
    } else if (level === 'warn') {
      Sentry.captureMessage(message, {
        extra: context,
        level: 'warning',
      });

      console.warn(JSON.stringify(entry));
    }
  }
}
