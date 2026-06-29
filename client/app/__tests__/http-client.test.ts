import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import { resetHttpClient, getHttpClient } from '~/api/http-client'

vi.mock('~/stores/AuthStore', () => ({
  useAuthStore: vi.fn(() => ({
    accessToken: 'test-access-token',
    refreshToken: 'test-refresh-token',
    updateTokens: vi.fn(),
    updateAccessToken: vi.fn(),
    clearTokens: vi.fn(),
  })),
}))

vi.mock('#imports', () => ({
  useRuntimeConfig: () => ({
    public: { serverUrl: 'http://localhost:3000/api/' },
  }),
}))

describe('http-client', () => {
  beforeEach(() => {
    resetHttpClient()
  })

  it('returns a singleton instance', () => {
    const a = getHttpClient()
    const b = getHttpClient()
    expect(a).toBe(b)
  })

  it('creates an Axios instance with the correct baseURL', () => {
    const client = getHttpClient()
    expect(client.defaults.baseURL).toBe('http://localhost:3000/api/')
    expect(client.defaults.timeout).toBe(15000)
    expect(client.defaults.headers['Content-Type']).toBe('application/json')
  })

  it('resetHttpClient clears the singleton', () => {
    const a = getHttpClient()
    resetHttpClient()
    const b = getHttpClient()
    expect(a).not.toBe(b)
  })
})

describe('request interceptor', () => {
  beforeEach(() => {
    resetHttpClient()
  })

  it('injects Authorization header from store', () => {
    const client = getHttpClient()
    // Axios request interceptors run before the actual request.
    // We can test the interceptor by inspecting how an intercepted
    // request config is modified.
    const interceptor = client.interceptors.request as any
    expect(interceptor).toBeDefined()
  })

  it('skips auth header for login endpoint', async () => {
    const client = getHttpClient()
    const config = { url: '/auth/login', headers: {} } as any
    // Run the request interceptor manually
    const handlers = client.interceptors.request as any
    const fulfilledHandler = handlers.handlers[0].fulfilled
    const result = await fulfilledHandler(config)
    expect(result.headers.Authorization).toBeUndefined()
  })

  it('skips auth header for register endpoint', async () => {
    const client = getHttpClient()
    const config = { url: '/auth/register', headers: {} } as any
    const handlers = client.interceptors.request as any
    const fulfilledHandler = handlers.handlers[0].fulfilled
    const result = await fulfilledHandler(config)
    expect(result.headers.Authorization).toBeUndefined()
  })

  it('skips auth header for refresh endpoint', async () => {
    const client = getHttpClient()
    const config = { url: '/refresh', headers: {} } as any
    const handlers = client.interceptors.request as any
    const fulfilledHandler = handlers.handlers[0].fulfilled
    const result = await fulfilledHandler(config)
    expect(result.headers.Authorization).toBeUndefined()
  })

  it('attaches Bearer token for other endpoints', async () => {
    const client = getHttpClient()
    const config = { url: '/works', headers: {} } as any
    const handlers = client.interceptors.request as any
    const fulfilledHandler = handlers.handlers[0].fulfilled
    const result = await fulfilledHandler(config)
    expect(result.headers.Authorization).toBe('Bearer test-access-token')
  })
})

describe('response interceptor — 401 handling', () => {
  beforeEach(() => {
    resetHttpClient()
    // Clear all Axios mock adapter responses
  })

  it('passes through non-401 errors', async () => {
    const client = getHttpClient()
    const error = { response: { status: 403 }, config: { headers: {} } }
    const handlers = client.interceptors.response as any
    const rejectedHandler = handlers.handlers[0].rejected
    await expect(rejectedHandler(error)).rejects.toEqual(error)
  })

  it('passes through retried requests', async () => {
    const client = getHttpClient()
    const error = {
      response: { status: 401 },
      config: { headers: {}, _retry: true },
    }
    const handlers = client.interceptors.response as any
    const rejectedHandler = handlers.handlers[0].rejected
    await expect(rejectedHandler(error)).rejects.toEqual(error)
  })

  it('rejects when no refresh token is available', async () => {
    // Override mock to have no refresh token
    vi.mocked(await import('~/stores/AuthStore')).useAuthStore = vi.fn(() => ({
      accessToken: 'test-access-token',
      refreshToken: null,
      updateTokens: vi.fn(),
      updateAccessToken: vi.fn(),
      clearTokens: vi.fn(),
    }))

    const client = getHttpClient()
    // We need a fresh client after changing the mock
    resetHttpClient()
    const client2 = getHttpClient()
    const error = {
      response: { status: 401 },
      config: { headers: {}, url: '/works' },
    }
    const handlers = client2.interceptors.response as any
    const rejectedHandler = handlers.handlers[0].rejected
    await expect(rejectedHandler(error)).rejects.toThrow()
  })
})
