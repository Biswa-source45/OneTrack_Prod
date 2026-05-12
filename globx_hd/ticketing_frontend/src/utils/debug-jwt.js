import { decodeJWT } from './jwt';

export function debugJWT(token) {
  const claims = decodeJWT(token);
  console.log('Decoded JWT claims:', claims);
  return claims;
}
