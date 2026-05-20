/**
 * Utility functions for JWT token handling
 */

interface JWTClaims {
    role?: string
    user_id?: number
    sub?: string
    exp?: number
    [key: string]: string | number | boolean | null | undefined
}

/**
 * Decode a JWT token without verification (client-side only)
 * @param token JWT token string
 * @returns Decoded payload or null if invalid
 */
export function decodeJWT(token: string,): JWTClaims | null {
    if (!token) return null;

    try {
    // JWT format: header.payload.signature
        const parts = token.split('.',);
        if (parts.length !== 3) return null;

        // Decode base64 URL encoded payload
        const payload = parts[1];
        const decoded = atob(payload.replace(/-/g, '+',).replace(/_/g, '/',),);
        return JSON.parse(decoded,);
    } catch (error) {
        console.error('Failed to decode JWT:', error,);
        return null;
    }
}

/**
 * Extract user role from JWT token
 * @param token JWT token string
 * @returns Role string or null if not found
 */
export function getRoleFromToken(token: string,): string | null {
    const decoded = decodeJWT(token,);
    if (!decoded) return null;

    // Check for role in standard claims
    if (decoded.role) return decoded.role;

    // Fallback to checking custom claims
    if (decoded['http://schemas.microsoft.com/ws/2008/06/identity/claims/role']) {
        return decoded['http://schemas.microsoft.com/ws/2008/06/identity/claims/role'];
    }

    return null;
}

/**
 * Extract user ID from JWT token
 * @param token JWT token string
 * @returns User ID number or null if not found
 */
export function getUserIdFromToken(token: string,): number | null {
    const decoded = decodeJWT(token,);
    if (!decoded) return null;

    // Check for user_id in standard claims
    if (decoded.user_id) return decoded.user_id;

    // Fallback to checking sub or nameidentifier
    if (decoded.sub) {
        const id = parseInt(decoded.sub, 10,);
        if (!isNaN(id,)) return id;
    }

    return null;
}

/**
 * Check if JWT token is expired
 * @param token JWT token string
 * @returns True if expired, false if valid or can't determine
 */
export function isTokenExpired(token: string,): boolean {
    const decoded = decodeJWT(token,);
    if (!decoded || !decoded.exp) return true;

    const expiryTime = decoded.exp * 1000; // Convert to milliseconds
    return Date.now() >= expiryTime;
}

/**
 * Get token expiration date
 * @param token JWT token string
 * @returns Date object or null if can't determine
 */
export function getTokenExpiry(token: string,): Date | null {
    const decoded = decodeJWT(token,);
    if (!decoded || !decoded.exp) return null;

    const expiryTime = decoded.exp * 1000; // Convert to milliseconds
    return new Date(expiryTime,);
}
