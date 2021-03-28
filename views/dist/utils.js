// bytesToSize 计算文件大小
function bytesToSize(bytes,num) {
    if (bytes === 0) return '0 B';
    var k = 1024;
    var sizes = ['B','KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var  i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(num) + ' ' + sizes[i];
}
// judgeFileType 判断文件类型
function judgeFileType(filename) {
    var reg = /(\\+)/g;
    var pfn = filename.replace(reg, "#");
    var arrpfn = pfn.split("#");
    var fn = arrpfn[arrpfn.length - 1];
    var arrfn = fn.split(".");
    var ext =arrfn[arrfn.length - 1];
    var imgExt=["jpg","jpeg","png","gif"];
    if (contains(imgExt, ext)){
        return "img"
    }
    return "other"
}

function contains(arr, obj) {
    var i = arr.length;
    while (i--) {
        if (arr[i] === obj) {
            return true;
        }
    }
    return false;
}