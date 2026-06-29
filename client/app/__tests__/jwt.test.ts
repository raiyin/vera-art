import { describe, it, expect } from 'vitest'
import {
  decodeJWT,
  getRoleFromToken,
  getUserIdFromToken,
  isTokenExpired,
  getTokenExpiry,
} from '~/utils/jwt'

function makeToken(payload: Record<string, unknown>): string {
  const b64 = btoa(JSON.stringify(payload))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '')
  return `header.${b64}.signature`
}

describe('decodeJWT', () => {
  it('returns null for empty string', () => {
    expect(decodeJWT('')).toBeNull()
  })

  it('returns null for malformed token (not 3 parts)', () => {
    expect(decodeJWT('a.b')).toBeNull()
  })

  it('decodes a valid token', () => {
    const token = makeToken({ sub: '42', role: 'admin', exp: 9999999999 })
    expect(decodeJWT(token)).toMatchObject({ sub: '42', role: 'admin', exp: 9999999999 })
  })

  it('returns null for invalid base64', () => {
    expect(decodeJWT('header.!!!.sig')).toBeNull()
  })
})

describe('getRoleFromToken', () => {
  it('returns role from standard claim', () => {
    expect(getRoleFromToken(makeToken({ role: 'admin' }))).toBe('admin')
  })

  it('returns null when no role claim', () => {
    expect(getRoleFromToken(makeToken({ sub: '1' }))).toBeNull()
  })

  it('returns null for invalid token', () => {
    expect(getRoleFromToken('bad')).toBeNull()
  })

  it('reads Microsoft-style role claim', () => {
    const token = makeToken({
      'http://schemas.microsoft.com/ws/2008/06/identity/claims/role': 'editor',
    })
    expect(getRoleFromToken(token)).toBe('editor')
  })
})

describe('getUserIdFromToken', () => {
  it('returns user_id from standard claim', () => {
    expect(getUserIdFromToken(makeToken({ user_id: 7 }))).toBe(7)
  })

  it('extracts id from sub claim', () => {
    expect(getUserIdFromToken(makeToken({ sub: '42' }))).toBe(42)
  })

  it('returns null for non-numeric sub', () => {
    expect(getUserIdFromToken(makeToken({ sub: 'abc' }))).toBeNull()
  })

  it('returns null for invalid token', () => {
    expect(getUserIdFromToken('bad')).toBeNull()
  })
})

describe('isTokenExpired', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })
  afterEach(() => {
    vi.useRealTimers()
  })

  it('returns true for token with no exp', () => {
    expect(isTokenExpired(makeToken({}))).toBe(true)
  })

  it('returns true for expired token', () => {
    vi.setSystemTime(new Date('2025-01-01'))
    const token = makeToken({ exp: 1735689600 }) // 2025-01-01 00:00:00 UTC
    vi.setSystemTime(new Date('2025-06-01'))
    expect(isTokenExpired(token)).toBe(true)
  })

  it('returns false for valid token', () => {
    vi.setSystemTime(new Date('2025-01-01'))
    const token = makeToken({ exp: 1767225600 }) // far future
    expect(isTokenExpired(token)).toBe(false)
  })

  it('returns true for invalid token', () => {
    expect(isTokenExpired('bad')).toBe(true)
  })
})

describe('getTokenExpiry', () => {
  it('returns Date for token with exp', () => {
    const token = makeToken({ exp: 1700000000 })
    const date = getTokenExpiry(token)
    expect(date).toBeInstanceOf(Date)
    expect(date!.getTime()).toBe(1700000000 * 1000)
  })

  it('returns null for token without exp', () => {
    expect(getTokenExpiry(makeToken({}))).toBeNull()
  })

  it('returns null for invalid token', () => {
    expect(getTokenExpiry('bad')).toBeNull()
  })
})
