import { Box, Button, Card, Flex, Grid, Loader, LoadingOverlay, Modal, Pagination, Table, Text, TextInput, Title } from "@mantine/core";
import { useDisclosure } from '@mantine/hooks';
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { QueryClient, useMutation, useQuery, useQueryClient } from "react-query";
import { useNavigate, useSearchParams } from "react-router-dom";
import { CreateAdminCourseRequest, FindAllAdminCoursesResponse_Course, PaginationMetadata } from "../../../pb/autograd/v1/autograd_pb";
import { AutogradServiceClient } from "../../../service";


function useListCourses(arg: {
    queryClient: QueryClient;
    page: number;
    limit: number;
}): {
    isLoading: boolean;
    isError: boolean;
    error: unknown;
    courses: FindAllAdminCoursesResponse_Course[];
    paginationMetadata: PaginationMetadata;
    isEmpty: () => boolean;
} {
    const queryKeys = ["courses", arg.page, arg.limit]

    const { isLoading, data, isError, error } = useQuery({
        queryKey: queryKeys,
        queryFn: async () => {
            return AutogradServiceClient.findAllAdminCourses({
                paginationRequest: {
                    page: arg.page,
                    limit: arg.limit,
                }
            })
        },
    })

    function isEmpty(): boolean {
        return !data || !data.courses || data?.courses.length === 0
    }

    return {
        isLoading,
        isError,
        error,
        courses: data?.courses || [],
        paginationMetadata: data?.paginationMetadata || new PaginationMetadata(),
        isEmpty,
    }
}

export function PageCourses() {
    const [overlayVisible, overlayMethod] = useDisclosure(false);

    const queryClient = useQueryClient()
    const [modalOpen, modalMethod] = useDisclosure(false);
    const navigate = useNavigate();

    const [searchParams] = useSearchParams()
    const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'))
    const limit = parseInt(searchParams.get('limit') || '10')
    const hookListCourses = useListCourses({
        queryClient,
        page,
        limit,
    })

    if (hookListCourses.isLoading) {
        return <>
            <Title order={3} mb="lg">
                Courses
            </Title>
            <Text>Loading...</Text>
        </>
    }

    if (hookListCourses.isError) {
        return <>
            <Title order={3} mb="lg">
                Courses
            </Title>
            <Text>Error: {hookListCourses.error as string}</Text>
        </>
    }

    function CourseHeading() {
        return <Flex direction="row">
            <Title order={3} mb="lg" mr="lg">
                Courses
            </Title>
        </Flex>
    }

    if (hookListCourses.isEmpty()) {
        return <>
            <CourseHeading />
            <Text><i>No courses</i></Text>
        </>
    }

    return <>
        <CourseHeading />

        <Box pos="relative">
            <LoadingOverlay visible={overlayVisible} zIndex={1000} overlayProps={{ radius: "sm", blur: 2 }} />
        </Box>

        <Grid>
            {
                hookListCourses.courses.map((course) => {
                    return <Grid.Col span={4} key={course.id}>
                        <Card
                            shadow="sm"
                            p="xl"
                            component="a"
                            target="_blank"
                            m="md"
                            style={{
                                '&:hover': {
                                    cursor: 'pointer'
                                }
                            }}
                            onClick={(e) => {
                                navigate(`/backoffice/courses/detail?courseID=${course.id}`)
                            }}
                        >
                            <Text fw={500} size="xl" mt="md">{course.name}</Text>
                            <Text mt="xs" c="dimmed" size="sm">{course.description}</Text>
                        </Card>
                    </Grid.Col>
                })
            }
        </Grid>

        <Pagination
            mb="lg"
            total={hookListCourses.paginationMetadata?.totalPage as number}
            value={page}
            onChange={setPage}
            siblings={1}
            boundaries={2}
        />
    </>
}