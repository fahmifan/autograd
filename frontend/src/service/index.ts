import {
	Interceptor,
	PromiseClient,
	createPromiseClient,
} from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { ServiceType } from "@bufbuild/protobuf";
import { useMemo } from "react";

import { AutogradService } from "../pb/autograd/v1/autograd_connect";

export function useAutogradClient(): PromiseClient<typeof AutogradService> {
	return useClient(AutogradService);
}

const csrfInterceptor: Interceptor = (next) => async (req) => {
	const token = localStorage.getItem("token")
	req.header.set("Authorization", `Bearer ${token}`);
	return await next(req);
};

const transport = createConnectTransport({
	baseUrl: "http://localhost:8080/grpc",
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
