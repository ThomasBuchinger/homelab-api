import { Card, CardContent, CardHeader, Typography } from "@mui/material"


export default function PlaceholderComponent(props: any) {
  var todos = props.todos || ["Nothing"]
  return <Card>
    <CardHeader title={props.name}></CardHeader>
      {todos.map((c: any) => {
        return <CardContent key={c}><Typography>{c}</Typography></CardContent>
      })}
    </Card>
}
