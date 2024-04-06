import { useQuery } from "react-query";
import type {
	FindAllStudentEnrolledCoursesResponse,
	FindStudentCourseDetailResponse,
} from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";

export function useCourseDetail(arg: {
	courseID: string;
}): {
	error: unknown;
	res?: FindStudentCourseDetailResponse;
} {
	const queryKeys = ["courses", arg.courseID];

	const { isLoading, data, isError, error } = useQuery({
		queryKey: queryKeys,
		queryFn: async () => {
			return AutogradServiceClient.findStudentCourseDetail({
				id: arg.courseID,
			});
		},
	});

	return {
		error,
		res: data,
	};
}

export function useListCourses(arg: {
	page: number;
	limit: number;
}): {
	isEmpty: () => boolean;
	isLoading: boolean;
	error: unknown;
	res?: FindAllStudentEnrolledCoursesResponse;
} {
	const queryKeys = ["student", "courses"];

	const { isLoading, data, isError, error } = useQuery({
		queryKey: queryKeys,
		queryFn: async () => {
			return AutogradServiceClient.findAllStudentEnrolledCourses({
				paginationRequest: {
					page: arg.page,
					limit: arg.limit,
				},
			});
		},
	});

	function isEmpty(): boolean {
		if (isLoading) {
			return true;
		}

		return (data?.courses?.length ?? 0) === 0;
	}

	return {
		isEmpty,
		isLoading,
		error,
		res: data,
	};
}
