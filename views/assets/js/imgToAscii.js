let toGrayScale = (r, g, b) => 0.299*r + 0.587*g + 0.114*b;
// const toGrayScale = (r, g, b) => 0.21 * r + 0.72 * g + 0.07 * b;

const grayRamp = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/|()1{}[]?-_+~<>i!lI;:,\"^`'. ";

// the grayScale value is an integer ranging from 0 (black) to 255 (white)
const grayScaleChar = grayScale => grayRamp[Math.ceil(((grayRamp.length - 1) * grayScale) / 255)];
let link = "assets/img/posts/whatsapp-logo.png"
async function imgToAscii(url) {
    let img = new Image()
    img.src = url
    let w, h;
    await new Promise((resolve) => {
        img.addEventListener('load', function () {
            [w, h] = [this.width, this.height]
            resolve();
        });
    });

    let canvas = document.createElement("canvas");
    let ctx = canvas.getContext("2d");
    canvas.width = w;
    canvas.height = h;

    ctx.drawImage(img, 0, 0);

    let rgba = ctx.getImageData(0, 0, img.width, img.height).data;
    let grayScales = [];

    for (let i = 0; i < rgba.length; i += 4) {
        let [r, g, b] = [rgba[i], rgba[i + 1], rgba[i + 2]];
        grayScales.push(toGrayScale(r, g, b))
    }

    let art = grayScales.map(c => grayScaleChar(c));

    let res = art.reduce((prev, next, i) => prev + next + ((i + 1) % w === 0 ? "\n" : ""), "")

    
    let pic = document.createElement("pre");
    pic.style.fontSize = "0.12vw";
    pic.style.color = "white";
    pic.textContent = res;
    pic.style.display = "inline-block";
    pic.style.transform = `scaleY(${h/w})`;

    document.querySelector("#all-posts .terminal-body").appendChild(pic);
    
   return res;
}


