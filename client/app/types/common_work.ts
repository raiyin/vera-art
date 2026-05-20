import type { TypedGetSaleDto, } from './sale';
import type { TypedGetWorkDto, } from './work';

export type CommonTypedGetWorkDto = TypedGetSaleDto | TypedGetWorkDto;
