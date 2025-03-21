import { RouteSectionProps } from "@solidjs/router";

export default function Layout(props: RouteSectionProps) {
  return <main>{props.children}</main>;
}
