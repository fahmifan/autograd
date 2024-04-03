import { Anchor, Breadcrumbs as MantineBreadcrumbs, BreadcrumbsFactory, BreadcrumbsProps, BreadcrumbsStylesNames, MantineStyleProps, StylesApiProps } from "@mantine/core";
import { Link } from "react-router-dom";

type Item = {
    title: string;
    to: string;
}

export function Breadcrumbs(props: {
    items: Item[];
}) {
    return <MantineBreadcrumbs>
        {
            props.items.map((item) => {
                return (
                    <Anchor key={item.to} component={Link} to={item.to}>
                        {item.title}
                    </Anchor>
                );
            })
        }
    </MantineBreadcrumbs>
}