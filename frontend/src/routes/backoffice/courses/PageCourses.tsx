import { Box, Button, Flex, Loader, LoadingOverlay, Modal, Pagination, Table, Text, TextInput, Title } from "@mantine/core";
import { useDisclosure } from '@mantine/hooks';
import { notifications } from "@mantine/notifications";
import { useState } from "react";
import { QueryClient, useMutation, useQuery, useQueryClient } from "react-query";
import { useSearchParams } from "react-router-dom";
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

    async function invalidateQueries(): Promise<void> {
        return arg.queryClient.invalidateQueries(queryKeys)
    }

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

function useCreateCourse(arg: {
    queryClient: QueryClient;
    onSuccess: (arg: { id: string }) => void;
    onError: (error: Error) => void;
}): {
    create(req: {
        name: string;
        description: string;
    }): Promise<void>;
    isLoading: boolean;
} {
    const [isLoading, setIsLoading] = useState(false)

    const mutation = useMutation({
        mutationFn: async (req: CreateAdminCourseRequest) => {
            return AutogradServiceClient.createAdminCourse({
                description: req.description,
                name: req.name,
            })
        },
        onError: (error: Error) => {
            arg.onError(error)
        },
    })

    async function create(req: {
        name: string;
        description: string;
    }): Promise<void> {
        setIsLoading(true)
        const res = await mutation.mutateAsync(new CreateAdminCourseRequest({
            description: req.description,
            name: req.name,
        }))
        await arg.queryClient.invalidateQueries(["courses"])
        arg.onSuccess({
            id: res.id
        })
        setIsLoading(false)
    }

    return {
        isLoading,
        create,
    }
}

export function PageCourses() {
    const [overlayVisible, overlayMethod] = useDisclosure(false);

    const queryClient = useQueryClient()
    const [modalOpen, modalMethod] = useDisclosure(false);

    const [searchParams] = useSearchParams()
    const [page, setPage] = useState(parseInt(searchParams.get('page') || '1'))
    const limit = parseInt(searchParams.get('limit') || '10')
    const hookListCourses = useListCourses({
        queryClient,
        page,
        limit,
    })

    const hookCreateCourse = useCreateCourse({
        queryClient,
        onSuccess: () => {
            overlayMethod.toggle()
            modalMethod.close()
        },
        onError: (error) => {
            notifications.show({
                message: error.message,
            })
        }
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
            <Button size="compact-md" onClick={() => {
                modalMethod.open()
            }}>Create</Button>

                <Modal opened={modalOpen} onClose={modalMethod.close} title="Authentication">
                    <form onSubmit={(e) => {
                        e.preventDefault()
                        const form = new FormData(e.target as HTMLFormElement)
                        hookCreateCourse.create({
                            name: form.get('name') as string,
                            description: form.get('description') as string,
                        })
                    }}>
                        <TextInput name="name" label="Course Name" placeholder="Name" />
                        <TextInput name="description" label="Course Description" placeholder="Description" />
                        
                        {
                            hookCreateCourse.isLoading 
                                ? <Loader color="blue" />
                                : <Button type="submit" mt="lg" disabled={hookCreateCourse.isLoading}>Save</Button>
                        }
                    </form>
                </Modal>
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

        <Table striped highlightOnHover mb="lg" maw={700}>
            <Table.Thead>
                <Table.Tr>
                    <Table.Th>Name</Table.Th>
                    <Table.Th>Description</Table.Th>
                    <Table.Th>Active</Table.Th>
                </Table.Tr>
            </Table.Thead>

            <Table.Tbody>
                {hookListCourses.courses.map((course) => {
                    return (
                        <Table.Tr key={course.id}>
                            <Table.Td>{course.name}</Table.Td>
                            <Table.Td>{course.description}</Table.Td>
                            <Table.Td>{course.active}</Table.Td>
                        </Table.Tr>
                    );
                })}
            </Table.Tbody>
        </Table>

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