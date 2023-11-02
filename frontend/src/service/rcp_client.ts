import { Err, OK, ResultFromPromise, ResultPromise } from "../types";

export class AutogradRPC {
	constructor(private baseUrl: string, private token: string) {
		this.baseUrl = baseUrl;
		this.token = token;
	}

	public async saveMedia(
		req: UploadMediaRequest,
	): ResultPromise<UploadMediaResponse> {
		// create form data
		const formData = new FormData();
		formData.append("media", req.file);
		formData.append("media_type", req.mediaType);

		const fetchres = await rfetch(`${this.baseUrl}/saveMedia`, {
			method: "POST",
			body: formData,
			headers: {
				Authorization: `Bearer ${this.token}`,
			},
		});
		if (!fetchres.ok) {
			return Err(fetchres.error);
		}

		if (fetchres.value.status < 200 || fetchres.value.status >= 400) {
			return Err(`Failed to save media with status ${fetchres.value.status}`);
		}

		const res = await (fetchres.value.json() as Promise<UploadMediaResponse>);

		return OK(res);
	}
}

type MediaFileType =
	| "assignment_case_input"
	| "assignment_case_output"
	| "submission";

export type UploadMediaResponse = {
	id: string;
};

export type UploadMediaRequest = {
	file: File;
	mediaType: MediaFileType;
};

async function rfetch(
	input: RequestInfo | URL,
	init?: RequestInit,
): ResultPromise<Response> {
	return ResultFromPromise(fetch(input, init));
}
