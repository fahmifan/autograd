import { useQuery } from "react-query";
import type { FindAdminCourseDetailResponse } from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";

export function useAdminCourseDetail(arg: {
	courseID: string;
}): {
	error: unknown;
	res?: FindAdminCourseDetailResponse;
} {
	const queryKeys = ["courses", arg.courseID];

	const { isLoading, data, isError, error } = useQuery({
		queryKey: queryKeys,
		queryFn: async () => {
			return AutogradServiceClient.findAdminCourseDetail({
				id: arg.courseID,
			});
		},
	});

	return {
		error,
		res: data,
	};
}
