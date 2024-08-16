import {
	Interceptor,
	createPromiseClient,
} from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";

import { JwtPayload, jwtDecode } from "jwt-decode";
import { AutogradQuery, AutogradService } from "../pb/autograd/v1/autograd_connect";
import { AutogradRPC } from "./rcp_client";


let _jwtToken = ''
export function getJWTToken(): string {
	if (_jwtToken) {
		return _jwtToken
	}
	_jwtToken = localStorage.getItem("token") ?? "";
	return _jwtToken
}

type JWTDecoded = JwtPayload & {
	id?: string;
	email?: string;
	name?: string;
	role?: string;
};

export function saveJWTToken(token: string): void {
	localStorage.setItem("token", token);
	_jwtToken = token
}

export function removeJWTToken(): void {
	localStorage.removeItem("token");
	_jwtToken = ''
}

export function decodeJWTToken(token: string): JWTDecoded {
	return jwtDecode<JWTDecoded>(token);
}

export function getDecodedJWTToken(): JWTDecoded {
	try {
		const token = getJWTToken()
		return decodeJWTToken(token)
	} catch(err) {
		return {}
	}
}

const csrfInterceptor: Interceptor = (next) => async (req) => {
	req.header.set("Authorization", `Bearer ${getJWTToken()}`);
	return await next(req);
};

const host = "http://localhost:8080";

const cmdTransport = createConnectTransport({
	baseUrl: `${host}/grpc`,
	interceptors: [csrfInterceptor],
});

const queryTransport = createConnectTransport({
	baseUrl: `${host}/grpc-query`,
	interceptors: [csrfInterceptor],
});

export const AutogradCmdClient = createPromiseClient(
	AutogradService,
	cmdTransport,
);

export const AutogradRPCCmdClient = new AutogradRPC(
	`${host}/api/v1/rpc`,
	getJWTToken(),
);

export const AutogradQueryClient = createPromiseClient(
	AutogradQuery,
	queryTransport,
)