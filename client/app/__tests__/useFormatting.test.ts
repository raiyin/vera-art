import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useFormatting } from '~/composables/useFormatting'

// Mock useI18n locale
vi.mock('#imports', () => ({
  useI18n: () => ({ locale: { value: 'ru' }, t: (k: string) => k }),
}))

describe('formatPrice', () => {
  const { formatPrice } = useFormatting()

  it('formats price in kopecks to rubles', () => {
    expect(formatPrice(150000)).toBe('1 500 ₽')
  })

  it('formats 0 as free', () => {
    expect(formatPrice(0)).toBe('Бесплатно')
  })

  it('uses a custom free label', () => {
    expect(formatPrice(0, { freeLabel: 'Free' })).toBe('Free')
  })

  it('returns fallback for null', () => {
    expect(formatPrice(null)).toBe('—')
  })

  it('returns fallback for undefined', () => {
    expect(formatPrice(undefined)).toBe('—')
  })

  it('handles price already in rubles', () => {
    expect(formatPrice(1500, { inKopecks: false })).toBe('1 500 ₽')
  })

  it('formats small prices correctly', () => {
    expect(formatPrice(99)).toBe('0,99 ₽')
  })

  it('formats large prices with thousand separators', () => {
    expect(formatPrice(10000000)).toBe('100 000 ₽')
  })
})

describe('formatDate', () => {
  const { formatDate } = useFormatting()

  it('formats an ISO date string', () => {
    const result = formatDate('2025-06-15')
    expect(result).toContain('2025')
    expect(result).toContain('июнь') // ru locale
  })

  it('returns fallback for null', () => {
    expect(formatDate(null)).toBe('—')
  })

  it('returns fallback for empty string', () => {
    expect(formatDate('')).toBe('—')
  })
})

describe('formatDateTime', () => {
  const { formatDateTime } = useFormatting()

  it('formats date and time', () => {
    const result = formatDateTime('2025-06-15T14:30:00')
    expect(result).toContain('2025')
    expect(result).toContain('14:30')
  })

  it('returns fallback for null', () => {
    expect(formatDateTime(null)).toBe('—')
  })
})

describe('formatShortDate', () => {
  const { formatShortDate } = useFormatting()

  it('formats short date', () => {
    const result = formatShortDate('2025-06-15')
    expect(result).toContain('15')
  })

  it('returns empty string for null', () => {
    expect(formatShortDate(null)).toBe('')
  })
})
