/**
 * 🛡️ Error Boundary Component
 * 
 * React runtime errors-ийг барьж app crash-аас хамгаална.
 * Production-д хэрэглэгчид ээлтэй error UI харуулна.
 * 
 * Features:
 * - ✅ Catches JavaScript errors anywhere in component tree
 * - ✅ Logs error to console for debugging
 * - ✅ Shows fallback UI instead of crashing
 * - ✅ Reload button for recovery
 * 
 * Usage:
 * ```tsx
 * <ErrorBoundary>
 *   <YourApp />
 * </ErrorBoundary>
 * ```
 * 
 * Note: Error boundaries do NOT catch errors for:
 * - Event handlers (use try-catch instead)
 * - Asynchronous code (setTimeout, requestAnimationFrame)
 * - Server-side rendering
 * - Errors thrown in the error boundary itself
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React from 'react'
import { Button } from './ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from './ui/card'

interface ErrorBoundaryProps {
  children: React.ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
  errorInfo: React.ErrorInfo | null
}

export class ErrorBoundary extends React.Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    }
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    // Update state so next render shows fallback UI
    return {
      hasError: true,
      error,
    }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Log error to console (can be sent to error reporting service)
    console.error('❌ ErrorBoundary caught an error:', error, errorInfo)
    
    this.setState({
      error,
      errorInfo,
    })

    // TODO: Send to error reporting service (Sentry, LogRocket, etc.)
    // Example: logErrorToService(error, errorInfo)
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    })
  }

  handleReload = () => {
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex min-h-screen items-center justify-center p-4 bg-background">
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle className="text-destructive">
                🚨 Алдаа гарлаа
              </CardTitle>
              <CardDescription>
                Уучлаарай, ямар нэг алдаа гарлаа. Та дахин оролдоно уу.
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-4">
              {/* Error message */}
              <div className="rounded-lg bg-destructive/10 p-4">
                <p className="text-sm font-mono text-destructive">
                  {this.state.error?.message || 'Тодорхойгүй алдаа'}
                </p>
              </div>

              {/* Technical details (development mode) */}
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <details className="text-xs">
                  <summary className="cursor-pointer font-semibold text-muted-foreground hover:text-foreground">
                    Техникийн дэлгэрэнгүй (dev only)
                  </summary>
                  <div className="mt-2 space-y-2 rounded bg-muted p-3 font-mono">
                    <div>
                      <strong>Error:</strong>
                      <pre className="mt-1 overflow-auto text-xs">
                        {this.state.error.toString()}
                      </pre>
                    </div>
                    {this.state.error.stack && (
                      <div>
                        <strong>Stack:</strong>
                        <pre className="mt-1 max-h-40 overflow-auto text-xs">
                          {this.state.error.stack}
                        </pre>
                      </div>
                    )}
                    {this.state.errorInfo && (
                      <div>
                        <strong>Component Stack:</strong>
                        <pre className="mt-1 max-h-40 overflow-auto text-xs">
                          {this.state.errorInfo.componentStack}
                        </pre>
                      </div>
                    )}
                  </div>
                </details>
              )}
            </CardContent>

            <CardFooter className="flex gap-2">
              <Button
                onClick={this.handleReset}
                variant="outline"
                className="flex-1"
              >
                Дахин оролдох
              </Button>
              <Button
                onClick={this.handleReload}
                className="flex-1"
              >
                Дахин ачаалах
              </Button>
            </CardFooter>
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}

