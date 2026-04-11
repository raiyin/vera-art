import { TypedGetSaleDto } from "./sale";
import { TypedGetWorkDto } from "./work";

export type CommonGetWorkDto = TypedGetSaleDto | TypedGetWorkDto;
