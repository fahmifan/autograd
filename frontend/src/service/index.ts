import {
	Interceptor,
	PromiseClient,
	createPromiseClient,
} from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { ServiceType } from "@bufbuild/protobuf";
import { useMemo } from "react";

import { JwtPayload, jwtDecode } from "jwt-decode";
import { AutogradService } from "../pb/autograd/v1/autograd_connect";
import { AutogradRPC } from "./rcp_client";

export function useAutogradClient(): PromiseClient<typeof AutogradService> {
	return useClient(AutogradService);
}

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

const transport = createConnectTransport({
	baseUrl: `${host}/grpc`,
	interceptors: [csrfInterceptor],
});

function useClient<T extends ServiceType>(service: T): PromiseClient<T> {
	// We memoize the client, so that we only create one instance per service.
	return useMemo(() => createPromiseClient(service, transport), [service]);
}

export const AutogradServiceClient = createPromiseClient(
	AutogradService,
	transport,
);

export const AutogradRPCClient = new AutogradRPC(
	`${host}/api/v1/rpc`,
	getJWTToken(),
);
