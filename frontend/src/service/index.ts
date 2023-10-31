import {
	Interceptor,
	PromiseClient,
	createPromiseClient,
} from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { ServiceType } from "@bufbuild/protobuf";
import { useMemo } from "react";

import { AutogradService } from "../pb/autograd/v1/autograd_connect";
import { AutogradRPC } from "./rcp_client";

export function useAutogradClient(): PromiseClient<typeof AutogradService> {
	return useClient(AutogradService);
}

const jwtToken = localStorage.getItem("token") ?? "";

const csrfInterceptor: Interceptor = (next) => async (req) => {
	req.header.set("Authorization", `Bearer ${jwtToken}`);
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
	jwtToken,
);
