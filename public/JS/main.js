//Comment
function showComment() {
    var commentArea = document.getElementById("comment-area");
    commentArea.classList.remove("hide");
}

// // upvote function
function upvote(Id) {
    fetch("http://localhost:8080/api/vote", {
        headers: {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
            "cache-control": "no-cache",
            "content-type": "application/x-www-form-urlencoded",
            "pragma": "no-cache",
            "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Windows\"",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-origin",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1"
        },
        referrer: "http://localhost:8080/post?id=" + Id,
        referrerPolicy: "strict-origin-when-cross-origin",
        body: "postId=" + Id + "&vote=1",
        method: "POST",
        mode: "cors",
        credentials: "include"
    }).then(() => {
        location.reload();
    });
}






// // Downvote function
function downvote(Id) {
    fetch("http://localhost:8080/api/vote", {
        headers: {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
            "cache-control": "no-cache",
            "content-type": "application/x-www-form-urlencoded",
            "pragma": "no-cache",
            "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Windows\"",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-origin",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1"
        },
        referrer: "http://localhost:8080/post?id=" + Id,
        referrerPolicy: "strict-origin-when-cross-origin",
        body: "postId=" + Id + "&vote=-1",
        method: "POST",
        mode: "cors",
        credentials: "include"
    }).then(() => {
        location.reload();
    });
}


// Upvote Comment Function
function upvoteComment(commentId) {
    fetch("http://localhost:8080/api/comments/votes", {
        headers: {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
            "cache-control": "no-cache",
            "content-type": "application/x-www-form-urlencoded",
            "pragma": "no-cache",
            "sec-ch-ua": "\"Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Windows\"",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-origin",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1"
        },
        referrer: "http://localhost:8080/comments/votes?id=" + commentId,
        referrerPolicy: "strict-origin-when-cross-origin",
        body: "commentId=" + commentId + "&vote=1",
        method: "POST",
        mode: "cors",
        credentials: "include"
    }).then(() => {
        location.reload();
    });
}

// Downvote Comment Function
function downvoteComment(commentId) {
    fetch("http://localhost:8080/api/comments/votes", {
        headers: {
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
            "cache-control": "no-cache",
            "content-type": "application/x-www-form-urlencoded",
            "pragma": "no-cache",
            "sec-ch-ua": "\"Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Windows\"",
            "sec-fetch-dest": "document",
            "sec-fetch-mode": "navigate",
            "sec-fetch-site": "same-origin",
            "sec-fetch-user": "?1",
            "upgrade-insecure-requests": "1"
        },
        referrer: "http://localhost:8080/comments/votes?id=" + commentId,
        referrerPolicy: "strict-origin-when-cross-origin",
        body: "commentId=" + commentId + "&vote=-1",
        method: "POST",
        mode: "cors",
        credentials: "include"
    }).then(() => {
        location.reload();
    });
}
