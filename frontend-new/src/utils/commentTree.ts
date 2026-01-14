import type { Comment } from "../types/models"

export interface CommentNode extends Comment { //Creates a subclass of Comment where each of them have an additonal attribute of children
    children: CommentNode[];
}

export const buildCommentTree = (comments: Comment[]): CommentNode[] =>{
    const commentMap : Map<number, CommentNode> = new Map(); //Mapping of ID: Comment
    const roots: CommentNode[] | undefined = []

    //Initialise all comments into the mapping
    comments.forEach((c) =>{
        commentMap.set(c.id, {...c,children: []})
    })
    comments.forEach((c)=>{
        if(!c.parent_comment_id.Valid){ //This means that it is a root comment
            let temp = commentMap.get(c.id)
            if(temp){
                roots.push(temp)
            }else{
                console.error("Error converting comments into tree")
            }
        }else{
            const p = commentMap.get(c.parent_comment_id.Int64) //Get the parent comment
            if(p){
                let temp = commentMap.get(c.id)
                !(temp) ? console.error("Error converting comments into tree") :
                p.children.push(temp) // Add it into the children attribute of the commentNode
            }else{ //Logically this will never get returned as all parent comment will be returned before the child
                let temp = commentMap.get(c.id)
                !(temp) ? console.error("Error converting comments into tree") :
                roots.push(temp)
            }
        }
    })
    return roots 
}