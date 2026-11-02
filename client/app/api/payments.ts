import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface CreatePaymentPayload {
    product_id: number
    promo_code?: string
}

export interface CreatePaymentResult {
    payment_id: number
    amount: number
    currency: string
    confirmation_url: string
    description: string
    status: string
}

export interface PaymentInfo {
    id: number
    user_id: number
    status: string
    amount: number
    currency: string
    description: string | null
    payment_method: string | null
    confirmation_url: string | null
    created_at: string
    updated_at: string
}

export interface PurchaseInfo {
    id: number
    user_id: number
    product_id: number
    payment_id: number
    price_paid: number
    status: string
    access_start: string
    access_end: string | null
    created_at: string
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Creates a new payment (redirects user to YooKassa).
 */
export async function createPayment(payload: CreatePaymentPayload,): Promise<CreatePaymentResult | null> {
    try {
        const { data, } = await getHttpClient().post<CreatePaymentResult>('payments/create', payload,);
        return data;
    } catch (e) {
        console.error('createPayment error:', e,);
        return null;
    }
}

/**
 * Gets payment status by ID.
 */
export async function getPaymentStatus(paymentId: number,): Promise<PaymentInfo | null> {
    try {
        const { data, } = await getHttpClient().get<PaymentInfo>(`payments/${paymentId}`,);
        return data;
    } catch (e) {
        console.error('getPaymentStatus error:', e,);
        return null;
    }
}

/**
 * Creates a purchase record after successful payment.
 */
export interface CreatePurchaseResult {
    purchase_id: number
    status: string
}

export async function createPurchase(paymentId: number,): Promise<CreatePurchaseResult | null> {
    try {
        const { data, } = await getHttpClient().post<CreatePurchaseResult>(
            'purchases',
            { payment_id: paymentId, },
        );
        return data;
    } catch (e) {
        console.error('createPurchase error:', e,);
        return null;
    }
}
