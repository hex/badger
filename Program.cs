// See https://aka.ms/new-console-template for more information

using System.Numerics;
using SixLabors.Fonts;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.Drawing;
using SixLabors.ImageSharp.Drawing.Processing;
using SixLabors.ImageSharp.PixelFormats;
using SixLabors.ImageSharp.Processing;
using Color = SixLabors.ImageSharp.Color;
using Path = System.IO.Path;
using Point = SixLabors.ImageSharp.Point;

const string outputDir = "badgerOutput";

void EmptyDirectory(string path)
{
    var di = new DirectoryInfo(path);

    foreach (var file in di.GetFiles()) file.Delete();
    foreach (var dir in di.GetDirectories()) dir.Delete(true);
}

string SanitizeFileName(string fileName)
{
    fileName = Path.GetFileNameWithoutExtension(fileName);
    fileName = fileName.Replace(" ", "_");
    return string.Join("_", fileName.Split(Path.GetInvalidFileNameChars()));
}

PointF GetPivotPoint(string pivot, Size badgeSize, Size imageSize)
{
    return pivot switch
    {
        "top" => new PointF((imageSize.Width - badgeSize.Width) / 2f, 0),
        "left" => new PointF(0, (imageSize.Height - badgeSize.Height) / 2f),
        "bottom" => new PointF((imageSize.Width - badgeSize.Width) / 2f, imageSize.Height - badgeSize.Height),
        "right" => new PointF(imageSize.Width - badgeSize.Width, (imageSize.Height - badgeSize.Height) / 2f),
        "topLeft" => new PointF(0, 0),
        "topRight" => new PointF(imageSize.Width - badgeSize.Width, 0),
        "bottomLeft" => new PointF(0, imageSize.Height - badgeSize.Height),
        "bottomRight" => new PointF(imageSize.Width - badgeSize.Width, imageSize.Height - badgeSize.Height),
        "center" => new PointF((imageSize.Width - badgeSize.Width) / 2f, (imageSize.Height - badgeSize.Height) / 2f),
        _ => new PointF((imageSize.Width - badgeSize.Width) / 2f, (imageSize.Height - badgeSize.Height) / 2f)
    };
}

TextAlignment GetTextAlignment(string alignment)
{
    return alignment switch
    {
        "left" => TextAlignment.Start,
        "center" => TextAlignment.Center,
        "right" => TextAlignment.End,
        _ => TextAlignment.Center
    };
}

void CreateBadge(string text, string icon, string fontName, int width, int height, string color,
    string textColor,
    string textAlignment,
    int angle, string badgePivot, int horizontalPadding, int verticalPadding, string horizontalPivot,
    string verticalPivot)

{
    var image = Image.Load(icon);
    var font = SystemFonts.CreateFont(fontName, 10, FontStyle.Bold);

    var blankImage = new Image<Rgba32>(image.Width, image.Height);

    var inputImageSize = image.Size();

    using var badgeImage = blankImage.Clone(ctx =>
    {
        var imageSize = ctx.GetCurrentSize();

        width = imageSize.Width * width / 100;
        height = imageSize.Height * height / 100;

        var pivot = GetPivotPoint(badgePivot, new Size(width, height), inputImageSize);

        var textHolder = new RectangularPolygon(pivot, new SizeF(width, height));
        var textSize = TextMeasurer.Measure(text, new TextOptions(font));
        var textScalingFactor = Math.Min(
            width / (textSize.Width + horizontalPadding / 2f),
            height / (textSize.Height + verticalPadding / 2f));
        var scaledFont = new Font(font, textScalingFactor * font.Size);

        var textOptions = new TextOptions(scaledFont)
        {
            Origin = new Vector2(textHolder.Center.X, textHolder.Center.Y),
            HorizontalAlignment = Enum.Parse<HorizontalAlignment>(horizontalPivot, true),
            VerticalAlignment = Enum.Parse<VerticalAlignment>(verticalPivot, true),
            TextAlignment = GetTextAlignment(textAlignment)
        };

        ctx.Fill(Color.Parse(color), textHolder);
        ctx.DrawText(textOptions, text, Color.Parse(textColor));
        ctx.Rotate(angle);
        ctx.Resize(inputImageSize);
    });


    badgeImage.Save($"{outputDir}{Path.DirectorySeparatorChar}{SanitizeFileName(icon)}_badge.png");
}

void AddBadge(string icon, int offsetX, int offsetY, float opacity, bool overwrite)
{
    var image = Image.Load(icon);
    var badgeName = $"{outputDir}{Path.DirectorySeparatorChar}{SanitizeFileName(icon)}_badge.png";
    image.Mutate(ctx =>
    {
        var badge = Image.Load(badgeName);
        ctx.DrawImage(badge, new Point(offsetX, offsetY), opacity);
    });

    Console.WriteLine(
        $"Writing to: {(overwrite ? icon : $"{outputDir}{Path.DirectorySeparatorChar}{Path.GetFileName(icon)}")}");

    image.Save(overwrite
        ? icon
        : $"{outputDir}{Path.DirectorySeparatorChar}{Path.GetFileName(icon)}");

    File.Delete($"{outputDir}{Path.DirectorySeparatorChar}{SanitizeFileName(icon)}_badge.png");
}

void Badger(
    [Option(null, "Set badge text")] string text,
    [Option(null, "Icon path.[.png | .jpg | .jpeg | .appiconset]")]
    string icon,
    [Option(null, "Font name")] string fontName = "Arial",
    [Option(null, "Badge width in percentage. 0 - 100 ")]
    int width = 100,
    [Option(null, "Badge height in percentage. 0 - 100 ")]
    int height = 20,
    [Option(null, "Set badge background color with a hexadecimal color code")]
    string color = "#4096EE",
    [Option(null, "Badge opacity")] float opacity = 1f,
    [Option(null, "Set badge text color with a hexadecimal color code")]
    string textColor = "#F9F7ED",
    [Option(null, "Set badge text alignment. left | center | right")]
    string textAlignment = "center",
    [Option("r", "Set badge rotation")] int angle = 0,
    [Option("x", "Set badge x-axis offset")]
    int offsetX = 0,
    [Option("y", "Set badge y-axis offset")]
    int offsetY = 0,
    [Option(null,
        "Set badge pivot point. top | left | bottom | right | topLeft | topRight | bottomLeft | bottomRight | center")]
    string badgePivot = "bottomLeft",
    [Option(null, "Set badge text horizontal padding")]
    int horizontalPadding = 5,
    [Option(null, "Set badge text vertical padding")]
    int verticalPadding = 0,
    [Option(null, "Set badge text horizontal pivot. left | center | right")]
    string horizontalPivot = "center",
    [Option(null, "Set badge text vertical pivot. top | center | bottom")]
    string verticalPivot = "center",
    [Option("o", "Replace input icon. WARNING: This will overwrite the input icon.")]
    bool overwrite = false
)
{
    Directory.CreateDirectory(outputDir);
    EmptyDirectory(outputDir);

    var fileAttributes = File.GetAttributes(icon);

    if (fileAttributes.HasFlag(FileAttributes.Directory))
    {
        Console.WriteLine($"Adding badge to all icons in {icon}");
        var files = Directory.EnumerateFiles(icon, "*.*", SearchOption.AllDirectories)
            .Where(s => s.EndsWith(".png") || s.EndsWith(".jpg") || s.EndsWith(".jpeg"));

        foreach (var file in files)
        {
            CreateBadge(text, file, fontName, width, height, color, textColor, textAlignment, angle,
                badgePivot,
                horizontalPadding,
                verticalPadding, horizontalPivot, verticalPivot);
            AddBadge(file, offsetX, offsetY, opacity, overwrite);
        }
    }
    else
    {
        CreateBadge(text, icon, fontName, width, height, color, textColor, textAlignment, angle, badgePivot,
            horizontalPadding,
            verticalPadding, horizontalPivot, verticalPivot);
        AddBadge(icon, offsetX, offsetY, opacity, overwrite);
    }

    if (overwrite)
        Directory.Delete(outputDir, true);
}

ConsoleApp.Run(args, Badger);